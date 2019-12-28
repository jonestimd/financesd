import TransactionModel, {ITransaction} from './TransactionModel';
import CategoryStore from '../store/CategoryStore';
import {computed, observable} from 'mobx';
import IMixedRowTableModel from './IMixedRowTableModel';

export default class TransactionTableModel implements IMixedRowTableModel<TransactionModel> {
    @observable
    transactions: TransactionModel[];
    static EMPTY: TransactionTableModel;

    constructor(transactions: ITransaction[], categoryStore: CategoryStore) {
        this.transactions = transactions.map(tx => new TransactionModel(tx, categoryStore)).sort(TransactionModel.compare);
        let balance = 0;
        for (const transaction of this.transactions) {
            transaction.balance = balance += transaction.subtotal;
        }
    }

    get groups() {
        return this.transactions;
    }

    getGroupIndex(rowIndex: number): [number, number] {
        let index = 0;
        let rows = 0;
        while (index < this.transactions.length) {
            const count = 1 + this.transactions[index].details.length;
            if (rows + count > rowIndex) break;
            rows += count;
            index += 1;
        }
        return [index, rows];
    }

    @computed
    get rowCount() {
        return this.transactions.reduce((count, transaction) => count + 1 + transaction.details.length, 0);
    }
}

TransactionTableModel.EMPTY = new TransactionTableModel([], null);