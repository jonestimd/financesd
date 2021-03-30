export interface IAsset {
    id: number;
    name: string;
    type: string;
    scale: number;
    symbol?: string;
    shares: number;
    firstAcquired?: string;
    costBasis?: number;
    dividends?: number;
    version: number;
    transactionCount: number;
}

export interface ISecurity extends IAsset {
    type: string;
}

export class SecurityModel implements ISecurity {
    id: number;
    name: string;
    type: string;
    scale: number;
    symbol?: string;
    shares: number;
    firstAcquired?: string;
    costBasis?: number;
    dividends?: number;
    version: number;
    transactionCount: number;

    constructor(security: ISecurity) {
        this.id = security.id;
        this.name = security.name;
        this.type = security.type;
        this.scale = security.scale;
        this.symbol = security.symbol;
        this.shares = security.shares;
        this.firstAcquired = security.firstAcquired;
        this.costBasis = security.costBasis;
        this.dividends = security.dividends;
        this.version = security.version;
        this.transactionCount = security.transactionCount;
    }

    get displayName() {
        return this.name + (this.symbol ? ` (${this.symbol})` : '');
    }
}
