import {computed} from 'mobx';

export type AmountType = 'DEBIT_DEPOSIT' | 'ASSET_VALUE';

export interface ICategory {
    id: string;
    code: string;
    description?: string;
    amountType: AmountType;
    parentId?: number;
    security: boolean;
    income: boolean;
    version: number;
    transactionCount: number;
}

interface ICategoryStore {
    getCategory: (id: string | number) => CategoryModel | undefined;
}

export class CategoryModel implements ICategory {
    id: string;
    code: string;
    description?: string;
    amountType: AmountType;
    parentId?: number;
    security: boolean;
    income: boolean;
    version: number;
    transactionCount: number;
    categoryStore: ICategoryStore;

    constructor(category: ICategory, categoryStore: ICategoryStore) {
        this.id = category.id;
        this.code = category.code;
        this.description = category.description;
        this.amountType = category.amountType;
        this.parentId = category.parentId;
        this.security = category.security;
        this.income = category.income;
        this.version = category.version;
        this.transactionCount = category.transactionCount;
        this.categoryStore = categoryStore;
    }

    @computed
    get parent() {
        return this.parentId ? this.categoryStore.getCategory(this.parentId) : undefined;
    }

    get isAssetValue() {
        return this.amountType === 'ASSET_VALUE';
    }

    @computed
    get displayName(): string {
        return this.parent ? `${this.parent.displayName}\u25BA${this.code}` : this.code;
    }
}
