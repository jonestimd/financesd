import TransactionModel, {ITransaction} from 'lib/model/TransactionModel';
import AccountStore from 'lib/store/AccountStore';
import CategoryStore from 'lib/store/CategoryStore';
import {RootStore} from 'lib/store/RootStore';

let nextId = 0;

export function newTx(overrides: Partial<ITransaction> = {}): ITransaction {
    return {
        id: ++nextId,
        version: 0,
        date: '2020-01-01',
        details: [],
        cleared: false,
        ...overrides,
    };
}

interface ModelOverrides extends Partial<ITransaction> {
    accountStore?: AccountStore;
    categoryStore?: CategoryStore;
    balance?: number;
}

const rootStore = new RootStore();

export function newTxModel({accountStore, categoryStore, balance, ...overrides}: ModelOverrides = {}): TransactionModel {
    const model = new TransactionModel(newTx(overrides), accountStore ?? rootStore.accountStore, categoryStore ?? rootStore.categoryStore);
    if (balance !== undefined) model.balance = balance;
    return model;
}
