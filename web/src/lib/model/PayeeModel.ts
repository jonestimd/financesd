export interface IPayee {
    id: number;
    name: string;
    version: number;
    transactionCount: number;
}

export class PayeeModel implements IPayee {
    id: number;
    name: string;
    version: number;
    transactionCount: number;

    constructor(payee: IPayee) {
        this.id = payee.id;
        this.name = payee.name;
        this.version = payee.version;
        this.transactionCount = payee.transactionCount;
    }
}
