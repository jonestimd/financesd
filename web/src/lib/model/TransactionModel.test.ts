import {newCategoryModel} from 'test/categoryFactory';
import {newDetail} from 'test/detailFactory';
import {newTx, newTxModel} from 'test/transactionFactory';
import {RootStore} from '../store/RootStore';
import TransactionModel from './TransactionModel';


describe('TransactionModel', () => {
    const {accountStore, categoryStore} = new RootStore();
    const transaction = newTx({details: [newDetail(), newDetail()]});
    const model = new TransactionModel(transaction, accountStore, categoryStore);

    beforeEach(() => model.reset());
    describe('constructor', () => {
        it('populates transaction properties', () => {
            const {details, ...tx} = transaction;
            expect(model).toEqual(expect.objectContaining(tx));
            details.forEach((detail, i) => expect(model.details[i]).toEqual(expect.objectContaining(detail)));
        });
    });
    describe('set date', () => {
        it('updates isChanged', () => {
            model.date = '0001-01-01';

            expect(model.isChanged).toBe(true);
        });
    });
    describe('set referenceNumber', () => {
        it('updates isChanged', () => {
            model.referenceNumber = '#123';

            expect(model.isChanged).toBe(true);
        });
    });
    describe('set payeeId', () => {
        it('updates isChanged', () => {
            model.payeeId = 123;

            expect(model.isChanged).toBe(true);
        });
    });
    describe('set securityId', () => {
        it('updates isChanged', () => {
            model.securityId = 123;

            expect(model.isChanged).toBe(true);
        });
    });
    describe('set memo', () => {
        it('updates isChanged', () => {
            model.memo = 'notes';

            expect(model.isChanged).toBe(true);
        });
    });
    describe('set cleared', () => {
        it('updates isChanged', () => {
            model.cleared = true;

            expect(model.isChanged).toBe(true);
        });
    });
    describe('get changes', () => {
        it('returns changes with id and version', () => {
            model.memo = 'notes';

            expect(model.changes).toEqual({memo: 'notes', id: model.id, version: model.version});
        });
    });
    describe('subtotal', () => {
        it('returns total of detail amounts', () => {
            expect(model.subtotal).toEqual(model.details[0].amount + model.details[1].amount);
        });
        it('ignores asset quantities', () => {
            const category = newCategoryModel({amountType: 'ASSET_VALUE'});
            categoryStore['categoriesById'].set(category.id, category);
            const details = [newDetail({amount: 12.34}), newDetail({transactionCategoryId: category.id})];
            const model = newTxModel({details, categoryStore});

            expect(model.subtotal).toEqual(details[0].amount);
        });
    });
    describe('getField', () => {
        const securityTests = [
            {transactionField: 'date', detailField: 'amount', itemIndex: -1},
            {transactionField: 'ref', detailField: 'category', itemIndex: -1},
            {transactionField: 'payee', detailField: 'group', itemIndex: -1},
            {transactionField: 'security', detailField: 'shares', itemIndex: -1},
            {transactionField: 'description', detailField: 'memo', itemIndex: -1},
            {transactionField: false, detailField: 'amount', itemIndex: 0},
            {transactionField: false, detailField: 'category', itemIndex: 0},
            {transactionField: false, detailField: 'group', itemIndex: 0},
            {transactionField: false, detailField: 'shares', itemIndex: 0},
            {transactionField: false, detailField: 'memo', itemIndex: 0},
        ];
        securityTests.forEach(({transactionField, detailField, itemIndex}, index) =>
            it(`returns ${transactionField} and ${detailField} for showSecurity and ${index}`, () => {
                expect(model.getField(index, true)).toEqual({transactionField, detailField, itemIndex});
            }));
        const notSecurityTests = [
            {transactionField: 'date', detailField: 'amount', itemIndex: -1},
            {transactionField: 'ref', detailField: 'category', itemIndex: -1},
            {transactionField: 'payee', detailField: 'group', itemIndex: -1},
            {transactionField: 'description', detailField: 'memo', itemIndex: -1},
            {transactionField: false, detailField: 'amount', itemIndex: 0},
            {transactionField: false, detailField: 'category', itemIndex: 0},
            {transactionField: false, detailField: 'group', itemIndex: 0},
            {transactionField: false, detailField: 'memo', itemIndex: 0},
        ];
        notSecurityTests.forEach(({transactionField, detailField, itemIndex}, index) =>
            it(`returns ${transactionField} and ${detailField} for !showSecurity and ${index}`, () => {
                expect(model.getField(index, false)).toEqual({transactionField, detailField, itemIndex});
            }));
    });
    describe('fieldCount', () => {
        it('returns 4 * (details + 1) for !showSecurity', () => {
            expect(model.fieldCount(false)).toEqual(4 * (model.details.length + 1));
        });
        it('returns 5 * (details + 1) for showSecurity', () => {
            expect(model.fieldCount(true)).toEqual(5 * (model.details.length + 1));
        });
    });
    describe('nextField', () => {
        const lastField = model.fieldCount(false)-1;

        it('returns 0 for 0', () => {
            expect(model.nextField(0, false)).toEqual(0);
        });
        it('returns last field for last field', () => {
            expect(model.nextField(lastField, false)).toEqual(lastField);
        });
        it('returns 0 for > last field', () => {
            expect(model.nextField(lastField+1, false)).toEqual(0);
        });
        it('returns last field for < 0', () => {
            expect(model.nextField(-1, false)).toEqual(lastField);
        });
    });
    describe('compare', () => {
        it('sorts by date', () => {
            const tx1 = newTx({date: '2020-12-01'});
            const tx2 = newTx({date: '2020-12-02'});

            expect(TransactionModel.compare(tx1, tx2)).toBeLessThan(0);
            expect(TransactionModel.compare(tx2, tx1)).toBeGreaterThan(0);
            expect(TransactionModel.compare(tx1, tx1)).toEqual(0);
            expect(TransactionModel.compare(tx2, tx2)).toEqual(0);
        });
        it('sorts undefined last', () => {
            const tx1 = newTx({date: undefined});
            const tx2 = newTx({date: '2020-12-02'});

            expect(TransactionModel.compare(tx1, tx2)).toBeGreaterThan(0);
            expect(TransactionModel.compare(tx2, tx1)).toBeLessThan(0);
            expect(TransactionModel.compare(tx1, tx1)).toEqual(0);
        });
    });
});
