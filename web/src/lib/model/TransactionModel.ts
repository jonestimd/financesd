import {makeObservable, observable} from 'mobx';
import CategoryStore from '../store/CategoryStore';

export interface IRelatedTransaction {
    id: string;
    accountId: number;
}

export interface IRelatedDetail {
    id: string;
    transaction: IRelatedTransaction;
}

export interface ITransactionDetail {
    id: string;
    transactionCategoryId: number;
    transactionGroupId: number;
    memo: string;
    amount: number;
    assetQuantity: number;
    relatedDetail: IRelatedDetail;
}

export interface ITransaction {
    id: string;
    date: string;
    referenceNumber?: string;
    payeeId?: number;
    securityId?: number;
    memo?: string;
    cleared: boolean;
    details: ITransactionDetail[];
}

export default class TransactionModel implements ITransaction {
    static compare(t1: ITransaction, t2: ITransaction): number {
        if (t1.date === undefined) return t2.date === undefined ? 0 : 1;
        if (t1.date === undefined) return -1;
        if (t1.date === t2.date) return Number(t1.id) - Number(t2.id);
        return t1.date < t2.date ? -1 : 1;
    }

    id: string;
    date: string;
    referenceNumber: string;
    payeeId: number;
    securityId: number;
    memo: string;
    cleared: boolean;
    @observable
    details: ITransactionDetail[];
    @observable
    subtotal: number;
    @observable
    balance: number;
    categoryStore: CategoryStore;

    constructor(transaction: ITransaction, categoryStore: CategoryStore) {
        makeObservable(this);
        Object.assign(this, transaction);
        this.categoryStore = categoryStore;
        this.subtotal = this.details.reduce((sum, detail) => this.isAssetValue(detail) ? sum : sum + detail.amount, 0);
    }

    private isAssetValue(detail: ITransactionDetail) {
        return this.categoryStore.getCategory(detail.transactionCategoryId).isAssetValue;
    }
}
