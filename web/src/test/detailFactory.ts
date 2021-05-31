import {ITransactionDetail} from 'lib/model/TransactionModel';

let nextId = 0;

export function newDetail(overrides: Partial<ITransactionDetail> = {}): ITransactionDetail {
    return {
        id: ++nextId,
        version: 0,
        amount: 123.78,
        ...overrides,
    };
}
