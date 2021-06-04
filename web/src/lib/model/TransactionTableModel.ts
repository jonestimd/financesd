import TransactionModel, {ITransaction, IUpdateTransactions} from './TransactionModel';
import AccountStore from 'lib/store/AccountStore';
import CategoryStore from '../store/CategoryStore';
import {action, computed, IReactionDisposer, makeObservable, observable, reaction, runInAction} from 'mobx';
import IMixedRowTableModel from './IMixedRowTableModel';
import {sortedIndex} from 'lodash';
import {INumberId} from './entityUtils';

export default class TransactionTableModel implements IMixedRowTableModel<TransactionModel> {
    readonly transactions = observable.array<TransactionModel>();
    private readonly accountStore: AccountStore;
    private readonly categoryStore: CategoryStore;
    private readonly balanceDisposer: IReactionDisposer;
    static EMPTY: TransactionTableModel;

    constructor(transactions: ITransaction[], accountStore: AccountStore, categoryStore: CategoryStore) {
        makeObservable(this);
        this.accountStore = accountStore;
        this.categoryStore = categoryStore;
        this.update(transactions);
        this.balanceDisposer = reaction(
            () => this.transactions.map((t) => t.subtotal),
            (subtotals) => runInAction(() => {
                subtotals.reduce((balance, subtotal, i) => {
                    return this.transactions[i].balance = subtotal + balance;
                }, 0);
            }),
            {fireImmediately: true},
        );
    }

    dispose() {
        this.balanceDisposer();
    }

    get groups() {
        return this.transactions;
    }

    @computed
    get precedingRows(): number[] {
        return this.transactions.reduce((rows, transaction) => {
            const last = rows[rows.length - 1];
            return rows.concat(last + 1 + transaction.details.length);
        }, [0]);
    }

    getGroupIndex(rowIndex: number): number {
        const index = sortedIndex(this.precedingRows, rowIndex);
        return this.precedingRows[index] === rowIndex ? index : index - 1;
    }

    getRowsAfter(groupIndex: number): number {
        return this.rowCount - this.precedingRows[groupIndex + 1];
    }

    @computed
    get rowCount() {
        return this.precedingRows[this.transactions.length];
    }

    @computed
    get isChanged() {
        return this.transactions.some((t) => t.isChanged);
    }

    @computed
    get isValid() {
        return ! this.transactions.some((t) => !t.isValid);
    }

    getTransaction(id: number) {
        return this.transactions.find((t) => t.id === id);
    }

    get changes(): IUpdateTransactions {
        const updates = this.transactions.filter((t) => t.isChanged).map((t) => t.changes);
        return {updates};
    }

    @action
    update(items: ITransaction[]) {
        const ids = items.map((t) => t.id);
        const transactions = this.transactions
            .filter((t) => !ids.includes(t.id))
            .concat(items.map((t) => new TransactionModel(t, this.accountStore, this.categoryStore)));
        this.transactions.replace(transactions.sort(TransactionModel.compare));
    }

    @action
    remove(items: INumberId[]) {
        const ids = items.map((item) => item.id);
        this.transactions.replace(this.transactions.filter((t) => !ids.includes(t.id)));
    }
}

TransactionTableModel.EMPTY = new TransactionTableModel([], {} as AccountStore, {} as CategoryStore);
