import TransactionModel, {ITransaction} from './TransactionModel';
import CategoryStore from '../store/CategoryStore';
import {computed, makeObservable, observable} from 'mobx';
import IMixedRowTableModel from './IMixedRowTableModel';
import {sortedIndex} from 'lodash';
import MessageStore from '../store/MessageStore';
import AlertStore from '../store/AlertStore';

export default class TransactionTableModel implements IMixedRowTableModel<TransactionModel> {
    @observable
    transactions: TransactionModel[];
    static EMPTY: TransactionTableModel;

    constructor(transactions: ITransaction[], categoryStore: CategoryStore) {
        makeObservable(this);
        this.transactions = transactions.map((tx) => new TransactionModel(tx, categoryStore)).sort(TransactionModel.compare);
        let balance = 0;
        for (const transaction of this.transactions) {
            transaction.balance = balance += transaction.subtotal;
        }
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
}

TransactionTableModel.EMPTY = new TransactionTableModel([], new CategoryStore(new MessageStore(), new AlertStore()));
