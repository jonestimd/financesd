import TransactionModel, {ITransaction} from 'src/lib/model/TransactionModel';
import CategoryStore from 'src/lib/store/CategoryStore';
import {defaultCategoryStore} from './categoryFactory';

let nextId = 0;

export function newTx(overrides: Partial<ITransaction> = {}): ITransaction {
    return {
        id: ++nextId,
        date: '2020-01-01',
        details: [],
        cleared: false,
        ...overrides,
    };
}

interface ModelOverrides extends Partial<ITransaction> {
    categoryStore?: CategoryStore;
    balance?: number;
}

export function newTxModel({categoryStore, balance, ...overrides}: ModelOverrides = {}): TransactionModel {
    const model = new TransactionModel(newTx(overrides), categoryStore ?? defaultCategoryStore);
    if (balance !== undefined) model.balance = balance;
    return model;
}
