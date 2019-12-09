export interface IPayee {
    id: string;
    name: string;
    version: number;
}

export class PayeeModel implements IPayee {
    id: string;
    name: string;
    version: number;

    constructor(payee: IPayee) {
        Object.assign(this, payee);
    }
}