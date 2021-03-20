import {newDetail} from 'src/test/detailFactory';
import {newTx} from 'src/test/transactionFactory';
import {RootStore} from '../store/RootStore';
import TransactionModel from './TransactionModel';
import TransactionTableModel from './TransactionTableModel';

describe('TransactionTableModel', () => {
    const {categoryStore} = new RootStore();
    const transactions = [
        newTx({details: [newDetail({amount: 12.34})]}),
        newTx({details: [newDetail({amount: 23.45})]}),
        newTx({details: [newDetail({amount: 34.56})]}),
    ];
    const model = new TransactionTableModel(transactions, categoryStore);

    describe('constructor', () => {
        it('creates TransactionModels', () => {
            expect(model.transactions[0]).toBeInstanceOf(TransactionModel);
        });
        it('calculates balance', () => {
            const subtotals = model.transactions.map((tx) => tx.subtotal);
            expect(model.transactions[0].balance).toEqual(subtotals[0]);
            expect(model.transactions[1].balance).toEqual(subtotals[0] + subtotals[1]);
            expect(model.transactions[2].balance).toEqual(subtotals[0] + subtotals[1] + subtotals[2]);
        });
    });
    describe('get groups', () => {
        it('returns the transactions', () => {
            expect(model.groups).toBe(model.transactions);
        });
    });
    describe('get precedingRows', () => {
        it('returns count of transactions+details', () => {
            const precedingRows = model.precedingRows;

            const detailCounts = transactions.map((tx) => tx.details.length);
            expect(precedingRows[0]).toEqual(0);
            expect(precedingRows[1]).toEqual(1 + detailCounts[0]);
            expect(precedingRows[2]).toEqual(2 + detailCounts[0] + detailCounts[1]);
            expect(precedingRows[3]).toEqual(3 + detailCounts[0] + detailCounts[1] + detailCounts[2]);
        });
    });
    describe('getGroupIndex', () => {
        it('returns transaction index for table row', () => {
            for (let i = 0; i < transactions.length * 2; i++) {
                expect(model.getGroupIndex(i)).toEqual(Math.floor(i / 2));
            }
        });
    });
    describe('getRowsAfter', () => {
        it('returns count of rows after transaction', () => {
            for (let i = 0; i < transactions.length; i++) {
                expect(model.getRowsAfter(i)).toEqual((transactions.length - i - 1) * 2);
            }
        });
    });
    describe('get rowCount', () => {
        it('returns count of transactions+details', () => {
            expect(model.rowCount).toEqual(transactions.length * 2);
        });
    });
});
