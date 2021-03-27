import {newCategoryModel} from 'src/test/categoryFactory';
import {newDetail} from 'src/test/detailFactory';
import {newTx, newTxModel} from 'src/test/transactionFactory';
import {RootStore} from '../store/RootStore';
import TransactionModel from './TransactionModel';


describe('TransactionModel', () => {
    const {categoryStore} = new RootStore();
    const transaction = newTx({details: [newDetail(), newDetail()]});

    describe('constructor', () => {
        it('populates transaction properties', () => {
            const model = new TransactionModel(transaction, categoryStore);

            expect(model).toEqual(expect.objectContaining(transaction));
        });
    });
    describe('subtotal', () => {
        it('returns total of detail amounts', () => {
            const model = new TransactionModel(transaction, categoryStore);

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
