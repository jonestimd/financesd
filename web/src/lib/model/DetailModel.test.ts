import {RootStore} from 'lib/store/RootStore';
import {newAccountModel} from 'test/accountFactory';
import {newCategoryModel} from 'test/categoryFactory';
import {newDetail, newDetailModel} from 'test/detailFactory';
import DetailModel from './DetailModel';

describe('DetailModel', () => {
    const {accountStore, categoryStore} = new RootStore();
    const detail = newDetail();
    const model = new DetailModel(accountStore, categoryStore, detail);
    const category = newCategoryModel();
    const account = newAccountModel();
    categoryStore['categoriesById'].set(category.id, category);
    accountStore['accountsById'].set(account.id, account);

    beforeEach(() => model.reset());
    describe('constructor', () => {
        it('sets empty values', () => {
            const model = new DetailModel(accountStore, categoryStore);

            expect(model.amount).toEqual(0);
            expect(model.amountText).toEqual('');
            expect(model.assetQuantityText).toEqual('');
            expect(model.isChanged).toBe(false);
            expect(model.isEmpty).toBe(true);
        });
        it('populates values with detail', () => {
            expect(model).toEqual(expect.objectContaining(detail));
            expect(model.amountText).toEqual(String(detail.amount));
            expect(model._assetQuantityText).toEqual('');
            expect(model.isChanged).toBe(false);
        });
        it('sets assetQuantityText', () => {
            const model = newDetailModel({assetQuantity: 123.456, accountStore, categoryStore});

            expect(model._assetQuantityText).toEqual('123.456');
        });
        it('sets transferAccountId', () => {
            const accountId = 42;

            const model = newDetailModel({relatedDetail: {id: -2, transaction: {id: -1, accountId}}, accountStore, categoryStore});

            expect(model.transferAccountId).toEqual(accountId);
        });
    });
    const setterTests = [
        {field: 'transactionCategoryId', value: 123, test: () => model.transactionCategoryId = 123},
        {field: 'transactionGroupId', value: 123, test: () => model.transactionGroupId = 123},
        {field: 'memo', value: 'notes', test: () => model.memo = 'notes'},
    ];
    setterTests.forEach(({field, value, test}) => {
        describe(`set ${field}`, () => {
            it('updates isChanged', () => {
                test();

                expect(model.isChanged).toBe(true);
                expect(model.changes).toEqual({id: model.id, version: model.version, [field]: value});
            });
        });
    });
    describe('get category', () => {
        it('returns a category', () => {
            const model = newDetailModel({categoryStore, transactionCategoryId: category.id});

            expect(model.category).toBe(category);
        });
        it('returns null for unknown category', () => {
            const model = newDetailModel({categoryStore, transactionCategoryId: -1});

            expect(model.category).toBe(null);
        });
        it('returns an account', () => {
            const model = newDetailModel({relatedDetail: {id: -2, transaction: {id: -1, accountId: account.id}}, accountStore, categoryStore});

            expect(model.category).toBe(account);
        });
        it('returns null for unknown account', () => {
            const model = newDetailModel({relatedDetail: {id: -2, transaction: {id: -1, accountId: -1}}, accountStore, categoryStore});

            expect(model.category).toBe(null);
        });
        it('returns null for no category or account', () => {
            const model = newDetailModel({accountStore, categoryStore});

            expect(model.category).toBe(null);
        });
    });
    describe('set category', () => {
        it('sets transactionCategoryId and clears transfer account', () => {
            const model = newDetailModel({relatedDetail: {id: -2, transaction: {id: -1, accountId: -1}}, accountStore, categoryStore});

            model.category = category;

            expect(model.changes).toEqual({id: model.id, version: model.version, transferAccountId: null, transactionCategoryId: category.id});
            expect(model.transactionCategoryId).toEqual(category.id);
            expect(model.transferAccountId).toBeUndefined();
        });
        it('sets transfer account and clears transactionCategoryId', () => {
            const model = newDetailModel({transactionCategoryId: category.id, accountStore, categoryStore});

            model.category = account;

            expect(model.changes).toEqual({id: model.id, version: model.version, transferAccountId: account.id, transactionCategoryId: null});
            expect(model.transactionCategoryId).toBeUndefined();
            expect(model.transferAccountId).toBe(account.id);
        });
        it('clears transactionCategoryId', () => {
            const model = newDetailModel({transactionCategoryId: category.id, accountStore, categoryStore});

            model.category = null;

            expect(model.changes).toEqual({id: model.id, version: model.version, transactionCategoryId: null});
            expect(model.transactionCategoryId).toBeUndefined();
            expect(model.transferAccountId).toBeUndefined();
        });
        it('clears transfer account', () => {
            const model = newDetailModel({relatedDetail: {id: -2, transaction: {id: -1, accountId: -1}}, accountStore, categoryStore});

            model.category = null;

            expect(model.changes).toEqual({id: model.id, version: model.version, transferAccountId: null});
            expect(model.transactionCategoryId).toBeUndefined();
            expect(model.transferAccountId).toBeUndefined();
        });
    });
    describe('get amount', () => {
        it('returns 0 for invalid amountText', () => {
            model.amountText = '';
            expect(model.amount).toEqual(0);

            model.amountText = '-';
            expect(model.amount).toEqual(0);
        });
    });
    describe('set amount', () => {
        it('sets amountText', () => {
            const model = newDetailModel({amount: 123.45});

            model.amount = 234.56;

            expect(model.amountText).toEqual('234.56');
            expect(model.amount).toEqual(234.56);
            expect(model.changes).toEqual({id: model.id, version: model.version, amount: 234.56});
        });
    });
    describe('set amountText', () => {
        it('sets amount', () => {
            const model = newDetailModel({amount: 123.45});

            model.amountText = '234.56';

            expect(model.amountText).toEqual('234.56');
            expect(model.amount).toEqual(234.56);
            expect(model.changes).toEqual({id: model.id, version: model.version, amount: 234.56});
        });
    });
    describe('set assetQuantity', () => {
        it('sets assetQuantityText', () => {
            const model = newDetailModel({assetQuantity: 123.456});

            model.assetQuantity = 234.567;

            expect(model.assetQuantityText).toEqual('234.567');
            expect(model.assetQuantity).toEqual(234.567);
            expect(model.changes).toEqual({id: model.id, version: model.version, assetQuantity: 234.567});
        });
    });
    describe('set assetQuantityText', () => {
        it('sets assetQuantity', () => {
            const model = newDetailModel({assetQuantity: 123.456});

            model.assetQuantityText = '234.567';

            expect(model.assetQuantityText).toEqual('234.567');
            expect(model.assetQuantity).toEqual(234.567);
            expect(model.changes).toEqual({id: model.id, version: model.version, assetQuantity: 234.567});
        });
        it('clears assetQuantity', () => {
            const model = newDetailModel({assetQuantity: 123.456});

            model.assetQuantityText = '';

            expect(model.assetQuantityText).toEqual('');
            expect(model.assetQuantity).toBeUndefined();
            expect(model.changes).toEqual({id: model.id, version: model.version, assetQuantity: null});
        });
    });
    describe('isChanged', () => {
        it('returns false if id exists', () => {
            model.amountText = '';

            expect(model.isEmpty).toBe(false);
        });
    });
    describe('isValid', () => {
        it('returns false for blank amount', () => {
            expect(newDetailModel({amountText: ''}).isValid).toBe(false);
        });
        it('returns false for invalid amount', () => {
            expect(newDetailModel({amountText: '-'}).isValid).toBe(false);
        });
        it('returns true for valid amount', () => {
            expect(newDetailModel({amountText: '0'}).isValid).toBe(true);
        });
    });
});
