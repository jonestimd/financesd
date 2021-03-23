import {ITransaction} from '../model/TransactionModel';
import {ObservableMap} from 'mobx';
import {RootStore} from './RootStore';
import TransactionTableModel from '../model/TransactionTableModel';
import Loader from './Loader';

export const query = `query($accountId: ID!) {
    transactions(accountId: $accountId) {
        id date referenceNumber payeeId securityId memo cleared
        details {
            id transactionCategoryId transactionGroupId memo amount assetQuantity
            relatedDetail {transaction {id accountId}}
        }
    }
}`;

export const loadingTransactions = 'Loading transactions';

export default class TransactionStore {
    private pendingAccounts: string[] = [];
    private transactionsByAccountId = new ObservableMap<string, TransactionTableModel>();
    private rootStore: RootStore;
    private loader: Loader;

    constructor(rootStore: RootStore) {
        // makeObservable(this);
        this.rootStore = rootStore;
        this.loader = new Loader(rootStore.messageStore, rootStore.alertStore);
    }

    getTransactionsModel(accountId?: string): TransactionTableModel {
        return accountId && this.transactionsByAccountId.get(accountId) || TransactionTableModel.EMPTY;
    }

    loadTransactions(accountId: string): Promise<boolean> | undefined {
        if (!this.transactionsByAccountId.has(accountId) && this.pendingAccounts.indexOf(accountId) < 0) {
            this.pendingAccounts.push(accountId);
            return this.loader.load<{transactions: ITransaction[]}>(loadingTransactions, {query, variables: {accountId},
                updater: ({transactions}) => {
                    this.transactionsByAccountId.set(accountId, new TransactionTableModel(transactions, this.rootStore.categoryStore));
                },
                completer: () => this.pendingAccounts.splice(this.pendingAccounts.indexOf(accountId), 1),
            });
        }
    }
}
