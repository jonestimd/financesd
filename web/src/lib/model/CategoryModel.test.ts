import {newCategory, newCategoryModel} from "src/test/categoryFactory";
import {RootStore} from "../store/RootStore";
import {CategoryModel} from "./CategoryModel";


describe('CategoryModel', () => {
    const {categoryStore} = new RootStore();
    const parent = newCategoryModel({categoryStore});
    const category = newCategory({parentId: parseInt(parent.id)});
    categoryStore['categoriesById'].set(parent.id, parent);

    describe('constructor', () => {
        it('populates category properties', () => {
            const model = new CategoryModel(category, categoryStore);

            expect(model).toEqual(expect.objectContaining(category));
        });
    });
    describe('get parent', () => {
        it('returns undefined for no parent', () => {
            expect(parent.parent).toBeUndefined();
        });
        it('returns parent model', () => {
            expect(new CategoryModel(category, categoryStore).parent).toBe(parent);
        });
    });
    describe('get isAssetValue', () => {
        it('returns true if amount type is ASSET_VALUE', () => {
            expect(newCategoryModel({amountType: 'ASSET_VALUE'}).isAssetValue).toBe(true);
        });
        it('returns false if amount type is not ASSET_VALUE', () => {
            expect(newCategoryModel({amountType: 'DEBIT_DEPOSIT'}).isAssetValue).toBe(false);
        });
    });
    describe('get displayName', () => {
        it('returns code if no parent', () => {
            expect(parent.displayName).toEqual(parent.code);
        });
        it('returns parent code and code', () => {
            expect(new CategoryModel(category, categoryStore).displayName).toEqual(`${parent.code}\u25ba${category.code}`);
        });
    });
});
