import * as agent from '../agent';
import {ITransaction} from '../model/TransactionModel';
import {flow, ObservableMap} from 'mobx';
import {RootStore} from './RootStore';
import TransactionTableModel from '../model/TransactionTableModel';
import {LoadResult} from './interfaces';

export const query = `query($accountId: ID!) {
    transactions(accountId: $accountId) {
        id date referenceNumber payeeId securityId memo cleared
        details {
            id transactionCategoryId transactionGroupId memo amount assetQuantity
            relatedDetail {transaction {id accountId}}
        }
    }
}`;

type TransactionsResponse = agent.IGraphqlResponse<{transactions: ITransaction[]}>;

export const loadingTransactions = 'Loading transactions...';

export default class TransactionStore {
    private pendingAccounts: string[] = [];
    private transactionsByAccountId = new ObservableMap<string, TransactionTableModel>();
    private rootStore: RootStore;

    constructor(rootStore: RootStore) {
        // makeObservable(this);
        this.rootStore = rootStore;
    }

    getTransactionsModel(accountId: string): TransactionTableModel {
        return this.transactionsByAccountId.get(accountId) || TransactionTableModel.EMPTY;
    }

    async loadTransactions(accountId: string): Promise<void> {
        if (!this.transactionsByAccountId.has(accountId) && this.pendingAccounts.indexOf(accountId) < 0) {
            this.rootStore.messageStore.addProgressMessage(loadingTransactions);
            await this._loadTransactions(accountId);
        }
    }

    private _loadTransactions = flow(function* (this: TransactionStore, accountId: string): LoadResult<TransactionsResponse> {
        this.pendingAccounts.push(accountId);
        try {
            const variables = {accountId};
            const {data} = yield agent.graphql('/finances/api/v1/graphql', query, variables);
            this.transactionsByAccountId.set(accountId, new TransactionTableModel(data.transactions, this.rootStore.categoryStore));
        } catch (err) {
            console.error('error gettting transactions', err);
        } finally {
            this.pendingAccounts.splice(this.pendingAccounts.indexOf(accountId), 1);
            this.rootStore.messageStore.removeProgressMessage(loadingTransactions);
        }
    });
}
