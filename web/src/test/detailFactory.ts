import DetailModel, {ITransactionDetail} from 'lib/model/DetailModel';
import AccountStore from 'lib/store/AccountStore';
import CategoryStore from 'lib/store/CategoryStore';
import {RootStore} from 'lib/store/RootStore';

let nextId = 0;

export function newDetail(overrides: Partial<ITransactionDetail> = {}): ITransactionDetail {
    return {
        id: ++nextId,
        version: 0,
        amount: 123.78,
        ...overrides,
    };
}

interface IModelOverrides extends Partial<ITransactionDetail> {
    amountText?: string;
    assetQuantityText?: string;
    accountStore?: AccountStore;
    categoryStore?: CategoryStore;
}

const rootStore = new RootStore();

export function newDetailModel({amountText, assetQuantityText, accountStore, categoryStore, ...overrides}: IModelOverrides = {}) {
    const model = new DetailModel(accountStore ?? rootStore.accountStore, categoryStore ?? rootStore.categoryStore, newDetail(overrides));
    if (amountText !== undefined) model.amountText = amountText;
    if (assetQuantityText !== undefined) model.assetQuantityText = assetQuantityText;
    return model;
}
