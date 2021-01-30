export interface IPayee {
    id: string;
    name: string;
    version: number;
    transactionCount: number;
}

export class PayeeModel implements IPayee {
    id: string;
    name: string;
    version: number;
    transactionCount: number;

    constructor(payee: IPayee) {
        Object.assign(this, payee);
    }
}
