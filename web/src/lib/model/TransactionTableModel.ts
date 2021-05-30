import TransactionModel, {ITransaction, IUpdateTransactions} from './TransactionModel';
import CategoryStore from '../store/CategoryStore';
import {action, autorun, computed, IReactionDisposer, makeObservable} from 'mobx';
import IMixedRowTableModel from './IMixedRowTableModel';
import {sortedIndex} from 'lodash';
import MessageStore from '../store/MessageStore';
import AlertStore from '../store/AlertStore';
import {INumberId} from './entityUtils';

export default class TransactionTableModel implements IMixedRowTableModel<TransactionModel> {
    transactions: TransactionModel[] = [];
    private readonly categoryStore: CategoryStore;
    private readonly balanceDisposer: IReactionDisposer;
    static EMPTY: TransactionTableModel;

    constructor(transactions: ITransaction[], categoryStore: CategoryStore) {
        makeObservable(this);
        this.categoryStore = categoryStore;
        this.update(transactions);
        this.balanceDisposer = autorun(() => {
            let balance = 0;
            for (const transaction of this.transactions) {
                transaction.balance = balance += transaction.subtotal;
            }
        });
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
            .concat(items.map((t) => new TransactionModel(t, this.categoryStore)));
        this.transactions = transactions.sort(TransactionModel.compare);
    }

    @action
    remove(items: INumberId[]) {
        const ids = items.map((item) => item.id);
        this.transactions = this.transactions.filter((t) => !ids.includes(t.id));
    }
}

TransactionTableModel.EMPTY = new TransactionTableModel([], new CategoryStore(new MessageStore(), new AlertStore()));
