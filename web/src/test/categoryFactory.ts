import {ICategory, CategoryModel} from 'lib/model/CategoryModel';
import AlertStore from 'lib/store/AlertStore';
import CategoryStore from 'lib/store/CategoryStore';
import MessageStore from 'lib/store/MessageStore';

let nextId = 0;

export function newCategory(overrides: Partial<ICategory> = {}): ICategory {
    return {
        id: ++nextId,
        code: `Category ${nextId}`,
        version: 1,
        amountType: 'DEBIT_DEPOSIT',
        income: false,
        security: false,
        transactionCount: 987,
        ...overrides,
    };
}

interface ModelOverrides extends Partial<ICategory> {
    categoryStore?: CategoryStore;
}

export const defaultCategoryStore = new CategoryStore(new MessageStore(), new AlertStore());

export function newCategoryModel({categoryStore, ...overrides}: ModelOverrides = {}): CategoryModel {
    return new CategoryModel(newCategory(overrides), categoryStore ?? defaultCategoryStore);
}
