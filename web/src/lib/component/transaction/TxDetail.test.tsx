import React from 'react';
import {shallow} from 'enzyme';
import TxDetail from './TxDetail';
import Category from './Category';
import {Currency, Shares} from 'src/lib/formats';
import Group from './Group';
import {newDetail} from 'src/test/detailFactory';

const detail = newDetail({transactionGroupId: 456});

describe('TxDetail', () => {
    it('shows amount, category and group', () => {
        const component = shallow(<TxDetail detail={detail} />);

        expect(component.find(Currency)).toHaveProp('children', detail.amount);
        expect(component.find(Category)).toHaveProp('detail', detail);
        expect(component.find(Group)).toHaveProp('id', detail.transactionGroupId);
        expect(component.find('.shares')).not.toExist();
        expect(component.find('.memo')).not.toExist();
    });
    it('shows asset quantity', () => {
        const assetQuantity = 100.789;

        const component = shallow(<TxDetail detail={{...detail, assetQuantity}} />);

        expect(component.find('.shares').find(Shares)).toHaveProp('children', assetQuantity);
    });
    it('shows memo', () => {
        const memo = 'note to self';

        const component = shallow(<TxDetail detail={{...detail, memo}} />);

        expect(component.find('.memo')).toHaveText(memo);
    });
});
