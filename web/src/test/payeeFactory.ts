import {IPayee, PayeeModel} from 'src/lib/model/PayeeModel';

let nextId = 0;

export function newPayee(overrides: Partial<IPayee> = {}): IPayee {
    return {
        id: `${++nextId}`,
        name: `Payee ${nextId}`,
        version: 1,
        transactionCount: 42,
        ...overrides,
    };
}

export function newPayeeModel(overrides: Partial<IPayee> = {}): PayeeModel {
    return new PayeeModel(newPayee(overrides));
}
