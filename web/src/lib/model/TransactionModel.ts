import {computed} from "mobx";

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
    balance: number;
    details: ITransactionDetail[];

    constructor(transaction: ITransaction) {
        Object.assign(this, transaction);
    }

    @computed
    get subtotal() {
        return this.details.reduce((sum, detail) => sum + detail.amount, 0);
    }

    static compare(t1: ITransaction, t2: ITransaction): number {
        if (t1.date === undefined) return t2.date === undefined ? 0 : 1;
        if (t1.date === undefined) return -1;
        if (t1.date === t2.date) return Number(t1.id) - Number(t2.id);
        return t1.date < t2.date ? -1 : 1;
    }
}