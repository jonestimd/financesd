import {computed, makeObservable, observable} from 'mobx';
import CategoryStore from '../store/CategoryStore';

export const securityTxFields = ['date', 'ref', 'payee', 'security', 'description'] as const;
export const txFields = [...securityTxFields].filter((name) => name !== 'security');
export type TransactionField = typeof securityTxFields[number];

export const securityDetailFields = ['amount', 'category', 'group', 'shares', 'memo'] as const;
export const detailFields = [...securityDetailFields].filter((name) => name !== 'shares');
export type DetailField = typeof securityDetailFields[number];

const fieldsPerItem = (showSecurity?: boolean) => showSecurity ? securityTxFields.length : txFields.length;

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
    @observable details: ITransactionDetail[];
    @observable balance = 0;
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

    getField(fieldIndex: number, showSecurity?: boolean) {
        const count = fieldsPerItem(showSecurity);
        const itemIndex = Math.floor(fieldIndex/count) - 1;
        const field = fieldIndex % count;
        const transactionField = itemIndex < 0 && (showSecurity ? securityTxFields[field] : txFields[field]);
        const detailField = showSecurity ? securityDetailFields[field] : detailFields[field];
        return {transactionField, detailField, itemIndex};
    }

    fieldCount(showSecurity?: boolean) {
        return fieldsPerItem(showSecurity) * (1 + this.details.length);
    }

    nextField(fieldIndex: number, showSecurity?: boolean) {
        if (fieldIndex < 0) return this.fieldCount(showSecurity) - 1;
        if (fieldIndex >= this.fieldCount(showSecurity)) return 0;
        return fieldIndex;
    }
}
