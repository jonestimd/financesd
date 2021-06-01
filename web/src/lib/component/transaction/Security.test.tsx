import React from 'react';
import {shallow} from 'enzyme';
import TransactionModel from 'lib/model/TransactionModel';
import {SecurityModel} from 'lib/model/SecurityModel';
import Security from './Security';
import {newTx} from 'test/transactionFactory';
import {RootStore} from 'lib/store/RootStore';
import {Icon} from '@material-ui/core';

describe('Security', () => {
    const rootStore = new RootStore();
    const {securityStore, accountStore, categoryStore} = rootStore;

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({securityStore});
    });
    it('returns null if no security', () => {
        const tx = new TransactionModel(newTx(), accountStore, categoryStore);

        const component = shallow(<Security transaction={tx} />);

        expect(component).toBeEmptyRender();
    });
    it('shows security', () => {
        const securityId = 123;
        const name = 'the security';
        jest.spyOn(securityStore, 'getSecurity').mockReturnValue({displayName: name} as SecurityModel);
        const tx = new TransactionModel(newTx({securityId}), accountStore, categoryStore);

        const component = shallow(<Security transaction={tx} />);

        expect(component).toHaveProp('data-type', 'security');
        expect(component).toHaveClassName('chip');
        expect(component.find(Icon)).toHaveText('request_page');
        expect(component.childAt(1)).toHaveText(name);
        expect(securityStore.getSecurity).toBeCalledWith(securityId);
    });
});
