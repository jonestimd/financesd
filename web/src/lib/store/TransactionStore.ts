import agent from 'superagent';
import {ITransaction} from '../model/TransactionModel';
import {flow, observable} from 'mobx';
import {RootStore} from './RootStore';
import TransactionTableModel from '../model/TransactionTableModel';

const query = `query($accountId: ID) {
    transactions(accountId: $accountId) {
        id date referenceNumber payeeId securityId memo cleared
        details {
            id transactionCategoryId transactionGroupId memo amount assetQuantity
            relatedDetail {transaction {id accountId}}
        }
    }
}`;

interface ITransactionsResponse {
    body: {data: {transactions: ITransaction[]}};
}

const loadingTransactions = 'Loading transactions...';

export default class TransactionStore {
    private pendingAccounts: string[] = [];
    @observable
    private transactionsByAccountId: {[accountId: string]: TransactionTableModel} = {};
    private rootStore: RootStore;

    constructor(rootStore: RootStore) {
        this.rootStore = rootStore;
    }

    getTransactionsModel(accountId: string): TransactionTableModel {
        return this.transactionsByAccountId[accountId] || TransactionTableModel.EMPTY;
    }

    loadTransactions(accountId: string): void {
        if (!this.transactionsByAccountId[accountId] && this.pendingAccounts.indexOf(accountId) < 0) {
            this.rootStore.messageStore.addProgressMessage(loadingTransactions);
            this._loadTransactions(accountId);
        }
    }

    private _loadTransactions = flow(function*(accountId: string) {
        this.pendingAccounts.push(accountId);
        try {
            const variables = {accountId};
            const {body: {data}}: ITransactionsResponse = yield agent.post('/finances/api/v1/graphql').send({query, variables});
            this.transactionsByAccountId[accountId] = new TransactionTableModel(data.transactions, this.rootStore.categoryStore);
        } catch (err) {
            console.error('error gettting transactions', err);
        } finally {
            this.pendingAccounts.splice(this.pendingAccounts.indexOf(accountId), 1);
            this.rootStore.messageStore.removeProgressMessage(loadingTransactions);
        }
    });
}