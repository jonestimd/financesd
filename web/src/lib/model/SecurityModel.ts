export interface IAsset {
    id: string;
    name: string;
    scale: number;
    symbol: string;
    version: number;
}

export interface ISecurity extends IAsset {
    type: string;
}

export class SecurityModel implements ISecurity {
    id: string;
    name: string;
    scale: number;
    symbol: string;
    version: number;
    type: string;

    constructor(security: ISecurity) {
        Object.assign(this, security);
    }
}