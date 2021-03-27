import {computed, makeObservable, observable} from 'mobx';
import CategoryStore from '../store/CategoryStore';

export interface IRelatedTransaction {
    id: number;
    accountId: number;
}

export interface IRelatedDetail {
    id: number;
    transaction: IRelatedTransaction;
}

export interface ITransactionDetail {
    id: number;
    transactionCategoryId?: number;
    transactionGroupId?: number;
    memo?: string;
    amount: number;
    assetQuantity?: number;
    relatedDetail?: IRelatedDetail;
}

export interface ITransaction {
    id: number;
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
        if (t2.date === undefined) return -1;
        if (t1.date === t2.date) return Number(t1.id) - Number(t2.id);
        return t1.date < t2.date ? -1 : 1;
    }

    id: number;
    date: string;
    referenceNumber?: string;
    payeeId?: number;
    securityId?: number;
    memo?: string;
    cleared: boolean;
    @observable
    details: ITransactionDetail[];
    @observable
    balance = 0;
    categoryStore: CategoryStore;

    constructor(transaction: ITransaction, categoryStore: CategoryStore) {
        makeObservable(this);
        this.id = transaction.id;
        this.date = transaction.date;
        this.referenceNumber = transaction.referenceNumber;
        this.payeeId = transaction.payeeId;
        this.securityId = transaction.securityId;
        this.memo = transaction.memo;
        this.cleared = transaction.cleared;
        this.details = transaction.details;
        this.categoryStore = categoryStore;
    }

    @computed
    get subtotal(): number {
        return this.details.reduce((sum, detail) => this.isAssetValue(detail) ? sum : sum + detail.amount, 0);
    }

    private isAssetValue({transactionCategoryId}: ITransactionDetail): boolean {
        return this.categoryStore.getCategory(transactionCategoryId)?.isAssetValue ?? false;
    }
}
