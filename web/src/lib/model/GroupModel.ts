export interface IGroup {
    id: string;
    name: string;
    description?: string;
    version: number;
    transactionCount: number;
}

export class GroupModel implements IGroup {
    id: string;
    name: string;
    description?: string;
    version: number;
    transactionCount: number;

    constructor(group: IGroup) {
        this.id = group.id;
        this.name = group.name;
        this.description = group.description;
        this.version = group.version;
        this.transactionCount = group.transactionCount;
    }
}
