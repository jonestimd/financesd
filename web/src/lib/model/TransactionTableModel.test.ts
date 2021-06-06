import {runInAction} from 'mobx';
import {newDetail} from 'test/detailFactory';
import {newTx} from 'test/transactionFactory';
import {RootStore} from '../store/RootStore';
import TransactionModel from './TransactionModel';
import TransactionTableModel from './TransactionTableModel';

describe('TransactionTableModel', () => {
    const {accountStore, categoryStore} = new RootStore();
    const transactions = [
        newTx({details: [newDetail({amount: 12.34})]}),
        newTx({details: [newDetail({amount: 23.45})]}),
        newTx({details: [newDetail({amount: 34.56})]}),
    ];
    let model: TransactionTableModel;
    beforeEach(() => model = new TransactionTableModel(transactions, accountStore, categoryStore));
    afterEach(() => model.dispose());
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
        it('updates balances', () => {
            runInAction(() => model.transactions[1].details[0].amount += 10);

            expect(model.transactions[0].balance).toBeCloseTo(12.34);
            expect(model.transactions[1].balance).toBeCloseTo(12.34 + 23.45 + 10);
            expect(model.transactions[2].balance).toBeCloseTo(12.34 + 23.45 + 10 + 34.56);
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
    describe('getTransaction', () => {
        it('return transaction with id', () => {
            const {details, ...rest} = transactions[0];
            expect(model.getTransaction(transactions[0].id)).toEqual(expect.objectContaining(rest));
        });
    });
    describe('isChanged', () => {
        it('returns false for no changes', () => {
            expect(model.isChanged).toBe(false);
        });
        it('returns true if a transaction has changes', () => {
            model.transactions[1].memo = 'notes';

            expect(model.isChanged).toBe(true);
        });
    });
    describe('get changes', () => {
        it('returns updates for changed transactions', () => {
            model.transactions[1].memo = 'notes';
            model.transactions[1].cleared = true;

            const changes = model.changes;

            const {id, version} = transactions[1];
            expect(changes).toEqual({updates: [{id, version, memo: 'notes', cleared: true, details: []}]});
        });
    });
    describe('update', () => {
        it('replaces transactions', () => {
            const {id} = transactions[1];
            const {details, ...updated} = newTx({id, version: transactions[1].version + 1, memo: 'notes'});

            model.update([{...updated, details}]);

            expect(model.getTransaction(id)).toEqual(expect.objectContaining(updated));
        });
        it('re-sorts transactions', () => {
            const {id} = transactions[1];
            const updated = newTx({id, version: transactions[1].version + 1, date: '1999-01-01'});

            model.update([updated]);

            expect(model.transactions[0].id).toEqual(id);
        });
    });
    describe('remove', () => {
        it('removes transactions with ids', () => {
            const {id} = transactions[1];

            model.remove([{id}]);

            expect(model.transactions.map((t) => t.id)).toEqual([transactions[0].id, transactions[2].id]);
        });
    });
});
