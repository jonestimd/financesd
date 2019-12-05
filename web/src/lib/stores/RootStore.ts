import React from 'react';

declare global {
    interface Window {
        rootStore: RootStore;
    }
}

export class RootStore {

}

export const RootStoreContext = React.createContext(new RootStore());
