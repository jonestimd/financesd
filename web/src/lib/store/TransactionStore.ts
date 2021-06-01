import {ITransaction} from '../model/TransactionModel';
import {ObservableMap} from 'mobx';
import {RootStore} from './RootStore';
import TransactionTableModel from '../model/TransactionTableModel';
import Loader from './Loader';

const transactionFields = `fragment transactionFields on transaction {
    id version date referenceNumber payeeId securityId memo cleared
    details {
        id version transactionCategoryId transactionGroupId memo amount assetQuantity
        relatedDetail {transaction {id accountId}}
    }
}`;

export const query = `${transactionFields}
query($accountId: Int!) {
    transactions(accountId: $accountId) {...transactionFields}
}`;

export const updateTxMutation = `${transactionFields}
mutation update($accountId: Int!, $deletes: [idVersion!], $adds: [addTransactionInput!], $updates: [updateTransactionInput!]) {
    transactions: updateTransactions(accountId: $accountId, delete: $deletes, add: $adds, update: $updates) {...transactionFields}
}`;

export const loadingTransactions = 'Loading transactions';
export const savingTransactions = 'Saving transactions';

export default class TransactionStore {
    private pendingAccounts: number[] = [];
    private transactionsByAccountId = new ObservableMap<number, TransactionTableModel>();
    private rootStore: RootStore;
    private loader: Loader;

    constructor(rootStore: RootStore) {
        // makeObservable(this);
        this.rootStore = rootStore;
        this.loader = new Loader(rootStore.messageStore, rootStore.alertStore);
    }

    getTransactionsModel(accountId?: number): TransactionTableModel {
        return typeof accountId === 'number' && this.transactionsByAccountId.get(accountId) || TransactionTableModel.EMPTY;
    }

    loadTransactions(accountId: number): Promise<boolean> | undefined {
        if (!this.transactionsByAccountId.has(accountId) && this.pendingAccounts.indexOf(accountId) < 0) {
            this.pendingAccounts.push(accountId);
            return this.loader.load<{transactions: ITransaction[]}>(loadingTransactions, {query, variables: {accountId},
                updater: ({transactions}) => {
                    const {accountStore, categoryStore} = this.rootStore;
                    this.transactionsByAccountId.set(accountId, new TransactionTableModel(transactions, accountStore, categoryStore));
                },
                completer: () => this.pendingAccounts.splice(this.pendingAccounts.indexOf(accountId), 1),
            });
        }
    }

    saveTransactions(accountId: number): Promise<boolean> {
        const tableModel = this.transactionsByAccountId.get(accountId);
        if (tableModel) {
            const changes = tableModel.changes;
            return this.loader.load<{transactions: ITransaction[]}>(savingTransactions, {query: updateTxMutation, variables: {accountId, ...changes},
                updater: ({transactions}) => {
                    // relatedAccountIds.forEach((id) => {
                    //     const model = this.transactionsByAccountId.get(id);
                    //     if (model) {
                    //         model.dispose();
                    //         this.transactionsByAccountId.delete(id);
                    //     }
                    // });
                    tableModel.update(transactions);
                    // if (changes.deletes) tableModel.remove(changes.deletes);
                },
            });
        }
        return Promise.resolve(true);
    }
}
