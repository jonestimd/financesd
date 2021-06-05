import {parseDate} from 'lib/formats';
import AccountStore from 'lib/store/AccountStore';
import {computed, makeObservable, observable} from 'mobx';
import CategoryStore from '../store/CategoryStore';
import ChangeModel from './ChangeModel';
import DetailModel, {IAddDetail, ITransactionDetail} from './DetailModel';
import {IVersionId, Nullable} from './entityUtils';

export const securityTxFields = ['date', 'ref', 'payee', 'security', 'description'] as const;
export const txFields = [...securityTxFields].filter((name) => name !== 'security');
export type TransactionField = typeof securityTxFields[number];

export const securityDetailFields = ['amount', 'category', 'group', 'shares', 'memo'] as const;
export const detailFields = [...securityDetailFields].filter((name) => name !== 'shares');
export type DetailField = typeof securityDetailFields[number];

const fieldsPerItem = (showSecurity?: boolean) => showSecurity ? securityTxFields.length : txFields.length;

export interface ITransaction {
    id: number;
    version: number;
    date: string;
    referenceNumber?: string;
    payeeId?: number;
    securityId?: number;
    memo?: string;
    cleared: boolean;
    details: ITransactionDetail[];
}

export interface IAddTransaction extends Omit<ITransaction, 'id' | 'version' | 'cleared' | 'details'> {
    cleared?: ITransaction['cleared'];
    details: IAddDetail[];
}

export interface IUpdateTransaction extends Nullable<Partial<Omit<ITransaction, 'id' | 'version' | 'details'>>> {
    id: ITransaction['id'];
    version: ITransaction['version'];
    details?: Partial<IAddDetail>[];
}

export interface IUpdateTransactions {
    adds?: IAddTransaction[];
    updates?: IUpdateTransaction[];
    deletes?: IVersionId[];
}

export default class TransactionModel implements ITransaction {
    static compare(t1: ITransaction, t2: ITransaction): number {
        if (t1.date === undefined) return t2.date === undefined ? 0 : 1;
        if (t2.date === undefined) return -1;
        if (t1.date === t2.date) return Number(t1.id) - Number(t2.id);
        return t1.date < t2.date ? -1 : 1;
    }

    id: number;
    @observable version: number;
    readonly details = observable.array<DetailModel>();
    @observable balance = 0;
    private _changes: ChangeModel<Omit<ITransaction, 'id' | 'version' | 'details'>>;
    // private readonly accountStore: AccountStore;
    private readonly categoryStore: CategoryStore;

    constructor(transaction: ITransaction, accountStore: AccountStore, categoryStore: CategoryStore) {
        makeObservable(this);
        this.id = transaction.id;
        this.version = transaction.version;
        this.details.replace(transaction.details.map((detail) => new DetailModel(detail, accountStore, categoryStore)));
        this._changes = new ChangeModel(transaction);
        // this.accountStore = accountStore;
        this.categoryStore = categoryStore;
    }

    @computed
    get date() {
        return this._changes.get('date');
    }

    set date(date: string) {
        this._changes.set('date', date);
    }

    @computed
    get referenceNumber() {
        return this._changes.get('referenceNumber');
    }

    set referenceNumber(ref: string | undefined) {
        this._changes.set('referenceNumber', ref ? ref : undefined);
    }

    @computed
    get payeeId() {
        return this._changes.get('payeeId');
    }

    set payeeId(payeeId: number | undefined) {
        this._changes.set('payeeId', payeeId);
    }

    @computed
    get securityId() {
        return this._changes.get('securityId');
    }

    set securityId(securityId: number | undefined) {
        this._changes.set('securityId', securityId);
    }

    @computed
    get memo() {
        return this._changes.get('memo');
    }

    set memo(memo: string | undefined) {
        this._changes.set('memo', memo ? memo : undefined);
    }

    @computed
    get cleared() {
        return this._changes.get('cleared');
    }

    set cleared(cleared: boolean) {
        this._changes.set('cleared', cleared);
    }

    get isChanged() {
        return this._changes.isChanged || this.details.some((d) => d.isChanged);
    }

    get changes(): IUpdateTransaction {
        return {...this._changes.changes, id: this.id, version: this.version};
    }

    reset() {
        this._changes.revert();
    }

    @computed
    get isValid() {
        return !!parseDate(this.date) && !this.details.some((d) => !d.isValid);
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
