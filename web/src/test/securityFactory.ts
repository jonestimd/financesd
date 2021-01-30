import {ISecurity, SecurityModel} from "src/lib/model/SecurityModel";

let nextId = 0;

export function newSecurity(overrides: Partial<ISecurity> = {}): ISecurity {
    return {
        id: `${++nextId}`,
        name: `Security ${nextId}`,
        type: 'Stock',
        scale: 6,
        shares: 0,
        version: 1,
        transactionCount: 96,
        ...overrides,
    };
}

export function newSecurityModel(overrides: Partial<ISecurity> = {}): SecurityModel {
    return new SecurityModel(newSecurity(overrides));
}
