import TransactionModel, {ITransaction} from './TransactionModel';
import CategoryStore from '../store/CategoryStore';
import {observable} from 'mobx';

export default class TransactionTableModel {
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
}

TransactionTableModel.EMPTY = new TransactionTableModel([], null);