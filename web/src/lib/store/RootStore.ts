import React from 'react';
import {action, observable} from 'mobx';
import {IMessageStore} from './MessageStore';
import {AccountStore} from './AccountStore';
import {CategoryStore} from './CategoryStore';
import {PayeeStore} from './PayeeStore';
import {TransactionStore} from './TransactionStore';

declare global {
    // tslint:disable-next-line: interface-name
    interface Window {
        rootStore: RootStore;
    }
}

export class RootStore implements IMessageStore {
    @observable
    private progressMessages: string[] = [];
    accountStore = new AccountStore(this);
    categoryStore = new CategoryStore(this);
    payeeStore = new PayeeStore(this);
    transactionStore = new TransactionStore(this);

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

export const RootStoreContext = React.createContext(new RootStore());
