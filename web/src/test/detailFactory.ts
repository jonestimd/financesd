import {ITransactionDetail} from 'src/lib/model/TransactionModel';

let nextId = 0;

export function newDetail(overrides: Partial<ITransactionDetail> = {}): ITransactionDetail {
    return {
        id: `${++nextId}`,
        amount: 123.78,
        ...overrides,
    };
}
