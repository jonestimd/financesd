import {computed} from "mobx";

export type AmountType = 'DEBIT_DEPOSIT' | 'ASSET_VALUE';

export interface ICategory {
    id: string;
    code: string;
    description: string;
    amountType: AmountType;
    parentId: number;
    security: boolean;
    income: boolean;
    version: number;
    transactionCount: number;
}

interface ICategoryStore {
    getCategory: (id: string | number) => CategoryModel;
}

export class CategoryModel implements ICategory {
    id: string;
    code: string;
    description: string;
    amountType: AmountType;
    parentId: number;
    security: boolean;
    income: boolean;
    version: number;
    transactionCount: number;
    categoryStore: ICategoryStore;

    constructor(category: ICategory, categoryStore: ICategoryStore) {
        Object.assign(this, category);
        this.categoryStore = categoryStore;
    }

    get parent() {
        return this.categoryStore.getCategory(this.parentId);
    }

    get isAssetValue() {
        return this.amountType === 'ASSET_VALUE';
    }

    @computed
    get displayName(): string {
        return this.parentId ? `${this.parent.displayName}\u25BA${this.code}` : this.code;
    }
}