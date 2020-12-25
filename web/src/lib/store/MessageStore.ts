import {action, makeObservable, observable} from 'mobx';

export interface IMessageStore {
    addProgressMessage(message: string): void;
    removeProgressMessage(message: string): void;
    readonly progressMessage: string | undefined;
}

export default class MessageStore implements IMessageStore {
    @observable
    private progressMessages: string[] = [];

    constructor() {
        makeObservable(this);
    }

    @action
    addProgressMessage(message: string) {
        this.progressMessages.push(message);
    }

    @action
    removeProgressMessage(message: string) {
        const index = this.progressMessages.indexOf(message);
        if (index >= 0) this.progressMessages.splice(index, 1);
    }

    get progressMessage(): string | undefined {
        return this.progressMessages.length > 0 ? this.progressMessages[0] : undefined;
    }
}
