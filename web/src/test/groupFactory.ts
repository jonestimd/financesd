import {IGroup, GroupModel} from 'lib/model/GroupModel';

let nextId = 0;

export function newGroup(overrides: Partial<IGroup> = {}): IGroup {
    return {
        id: ++nextId,
        name: `Group ${nextId}`,
        transactionCount: 0,
        version: 1,
        ...overrides,
    };
}

export function newGroupModel(overrides: Partial<IGroup> = {}): GroupModel {
    return new GroupModel(newGroup(overrides));
}
