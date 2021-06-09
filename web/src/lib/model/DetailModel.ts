import AccountStore from 'lib/store/AccountStore';
import CategoryStore from 'lib/store/CategoryStore';
import {action, computed, makeObservable, observable} from 'mobx';
import {AccountModel} from './account/AccountModel';
import {CategoryModel} from './CategoryModel';
import ChangeModel from './ChangeModel';
import {Nullable} from './entityUtils';

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
    version: number;
    transactionCategoryId?: number;
    transactionGroupId?: number;
    memo?: string;
    amount: number;
    assetQuantity?: number;
    relatedDetail?: IRelatedDetail;
}

export interface IAddDetail extends Omit<ITransactionDetail, 'id' | 'version' | 'relatedDetail'> {
    transferAccountId?: number
}

export type IUpdateDetail = IAddDetail | (Nullable<Partial<IAddDetail>> & Required<Pick<ITransactionDetail, 'id' | 'version'>>);

const toString = (value?: number) => typeof value === 'number' ? String(value) : '';

const isCategory = (x: CategoryModel | AccountModel): x is CategoryModel => 'code' in x;

export default class DetailModel {
    private readonly accountStore: AccountStore;
    private readonly categoryStore: CategoryStore;
    id?: number;
    @observable version?: number;
    @observable _amountText: string;
    @observable _assetQuantityText: string;
    @observable relatedDetail?: IRelatedDetail;
    private _changes: ChangeModel<IAddDetail & {amountText: string, assetQuantityText: string}>;

    constructor(accountStore: AccountStore, categoryStore: CategoryStore, detail?: ITransactionDetail) {
        makeObservable(this);
        this.accountStore = accountStore;
        this.categoryStore = categoryStore;
        this.id = detail?.id;
        this.version = detail?.version;
        this._amountText = detail ? String(detail.amount) : '';
        this._assetQuantityText = toString(detail?.assetQuantity);
        this.relatedDetail = detail?.relatedDetail;
        this._changes = new ChangeModel({
            ...(detail ?? {amount: NaN}),
            transferAccountId: detail?.relatedDetail?.transaction.accountId,
            amountText: this._amountText,
            assetQuantityText: this._assetQuantityText,
        });
    }

    @computed
    get category(): CategoryModel | AccountModel | null {
        if (this.transactionCategoryId) return this.categoryStore.getCategory(this.transactionCategoryId) ?? null;
        if (this.relatedDetail) return this.accountStore.getAccount(this.relatedDetail.transaction.accountId) ?? null;
        return null;
    }

    set category(categoryOrAccount: CategoryModel | AccountModel | null) {
        if (categoryOrAccount) {
            if (isCategory(categoryOrAccount)) {
                this.transactionCategoryId = categoryOrAccount.id;
                this.transferAccountId = undefined;
            }
            else {
                this.transferAccountId = categoryOrAccount.id;
                this.transactionCategoryId = undefined;
            }
        }
        else {
            this.transactionCategoryId = undefined;
            this.transferAccountId = undefined;
        }
    }

    @computed
    get transactionCategoryId() {
        return this._changes.get('transactionCategoryId');
    }

    set transactionCategoryId(categoryId: number | undefined) {
        this._changes.set('transactionCategoryId', categoryId);
    }

    @computed
    get transferAccountId() {
        return this._changes.get('transferAccountId');
    }

    set transferAccountId(accountId: number | undefined) {
        this._changes.set('transferAccountId', accountId);
    }

    @computed
    get transactionGroupId() {
        return this._changes.get('transactionGroupId');
    }

    set transactionGroupId(GroupId: number | undefined) {
        this._changes.set('transactionGroupId', GroupId);
    }

    @computed
    private get _amount() {
        return this._changes.get('amount');
    }

    @computed
    get amount() {
        return isNaN(this._amount) ? 0 : this._amount;
    }

    set amount(amount: number) {
        this._changes.set('amount', amount);
        this._amountText = String(amount);
    }

    @computed
    get amountText() {
        return this._amountText;
    }

    set amountText(amount: string) {
        this._amountText = amount;
        this._changes.set('amount', parseFloat(amount));
    }

    @computed
    get assetQuantity() {
        return this._changes.get('assetQuantity');
    }

    set assetQuantity(assetQuantity: number | undefined) {
        this._changes.set('assetQuantity', assetQuantity);
        this._assetQuantityText = toString(assetQuantity);
    }

    @computed
    get assetQuantityText() {
        return this._assetQuantityText;
    }

    set assetQuantityText(assetQuantity: string) {
        this._assetQuantityText = assetQuantity;
        this._changes.set('assetQuantity', assetQuantity ? parseFloat(assetQuantity) : undefined);
    }

    @computed
    get memo() {
        return this._changes.get('memo');
    }

    set memo(memo: string | undefined) {
        this._changes.set('memo', memo);
    }

    get isChanged() {
        return this._changes.isChanged;
    }

    get isValid() {
        return this.amountText != '' && !isNaN(this._amount);
    }

    get isEmpty() {
        return this.id === undefined && this.amountText === '';
    }

    get changes(): IUpdateDetail {
        const {amountText, assetQuantityText, ...changes} = this._changes.changes;
        return {...changes, id: this.id, version: this.version} as IAddDetail;
    }

    @action
    reset() {
        this._changes.revert();
    }
}
