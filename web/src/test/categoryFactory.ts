import {ICategory, CategoryModel} from "src/lib/model/CategoryModel";
import CategoryStore from "src/lib/store/CategoryStore";

let nextId = 0;

export function newCategory(overrides: Partial<ICategory> = {}): ICategory {
    return {
        id: `${++nextId}`,
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

export function newCategoryModel({categoryStore, ...overrides}: ModelOverrides = {}): CategoryModel {
    return new CategoryModel(newCategory(overrides), categoryStore);
}
