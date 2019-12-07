import React from 'react';
import {IMessageStore} from './MessageStore';
import {AccountStore} from './AccountStore';
import {action, observable} from 'mobx';

declare global {
    // tslint:disable-next-line: interface-name
    interface Window {
        rootStore: RootStore;
    }
}

export class RootStore implements IMessageStore {
    @observable
    private progressMessages: string[] = [];
    private _accountStore = new AccountStore(this);

    @action
    addProgressMessage(message: string) {
        this.progressMessages.push(message);
    }

    @action
    removeProgressMessage(message: string) {
        const index = this.progressMessages.indexOf(message);
        if (index >= 0) this.progressMessages.splice(index, 1);
    }

    get accountStore() {
        return this._accountStore;
    }

    get progressMessage(): string | undefined {
        return this.progressMessages.length > 0 ? this.progressMessages[0] : undefined;
    }
}

export const RootStoreContext = React.createContext(new RootStore());
