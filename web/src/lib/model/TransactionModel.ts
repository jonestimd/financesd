import {computed} from "mobx";
import {ICategoryStore} from '../store/CategoryStore';

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

export interface ITransaction extends IRelatedTransaction {
    date: string;
    referenceNumber: string;
    payeeId: number;
    memo: string;
    cleared: boolean;
    details: ITransactionDetail[];
}

export class TransactionModel implements ITransaction {
    id: string;
    accountId: number;
    date: string;
    referenceNumber: string;
    payeeId: number;
    memo: string;
    cleared: boolean;
    details: ITransactionDetail[];
    previous: TransactionModel;
    categoryStore: ICategoryStore;

    constructor(transaction: ITransaction, categoryStore: ICategoryStore) {
        Object.assign(this, transaction);
        this.categoryStore = categoryStore;
    }

    private isAssetValue(detail: ITransactionDetail) {
        return this.categoryStore.getCategory(detail.transactionCategoryId).isAssetValue;
    }

    @computed
    get subtotal() {
        return this.details.reduce((sum, detail) => this.isAssetValue(detail) ? sum : sum + detail.amount, 0);
    }

    @computed
    get balance(): number {
        const previousBalance = this.previous ? this.previous.balance : 0;
        return previousBalance + this.subtotal;
    }

    static compare(t1: ITransaction, t2: ITransaction): number {
        if (t1.date === undefined) return t2.date === undefined ? 0 : 1;
        if (t1.date === undefined) return -1;
        if (t1.date === t2.date) return Number(t1.id) - Number(t2.id);
        return t1.date < t2.date ? -1 : 1;
    }
}