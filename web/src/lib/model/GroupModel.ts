export interface IGroup {
    id: string;
    name: string;
    description?: string;
    version: number;
}

export class GroupModel implements IGroup {
    id: string;
    name: string;
    description?: string;
    version: number;

    constructor(group: IGroup) {
        Object.assign(this, group);
    }
}
