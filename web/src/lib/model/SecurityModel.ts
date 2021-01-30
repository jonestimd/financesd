export interface IAsset {
    id: string;
    name: string;
    type: string;
    scale: number;
    symbol?: string;
    version: number;
    transactionCount: number;
}

export interface ISecurity extends IAsset {
    type: string;
}

export class SecurityModel implements ISecurity {
    id: string;
    name: string;
    type: string;
    scale: number;
    symbol: string;
    version: number;
    transactionCount: number;

    constructor(security: ISecurity) {
        Object.assign(this, security);
    }
}
