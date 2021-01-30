import React from 'react';
import {shallow} from 'enzyme';
import TransactionModel from 'src/lib/model/TransactionModel';
import MessageStore from 'src/lib/store/MessageStore';
import SecurityStore from 'src/lib/store/SecurityStore';
import {SecurityModel} from 'src/lib/model/SecurityModel';
import Security from './Security';
import {newTx} from 'src/test/transactionFactory';
import CategoryStore from 'src/lib/store/CategoryStore';

describe('Security', () => {
    const messageStore = new MessageStore();
    const securityStore = new SecurityStore(messageStore);
    const categoryStore = new CategoryStore(messageStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({securityStore});
    });
    it('returns null if no security', () => {
        const tx = new TransactionModel(newTx(), categoryStore);

        const component = shallow(<Security transaction={tx} />);

        expect(component).toBeEmptyRender();
    });
    it('shows payee', () => {
        const securityId = 123;
        const name = 'the security';
        jest.spyOn(securityStore, 'getSecurity').mockReturnValue({name} as SecurityModel);
        const tx = new TransactionModel(newTx({securityId}), categoryStore);

        const component = shallow(<Security transaction={tx} />);

        expect(component).toHaveClassName('security');
        expect(component).toHaveClassName('chip');
        expect(component.find('i.material-icons.md-18')).toHaveText('request_page');
        expect(component.childAt(1)).toHaveText(name);
        expect(securityStore.getSecurity).toBeCalledWith(securityId);
    });
});
