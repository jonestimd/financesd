import React from 'react';
import {shallow} from 'enzyme';
import MessageStore from 'lib/store/MessageStore';
import GroupStore from 'lib/store/GroupStore';
import Group from './Group';
import {GroupModel} from 'lib/model/GroupModel';
import AlertStore from 'lib/store/AlertStore';

describe('Group', () => {
    const groupStore = new GroupStore(new MessageStore(), new AlertStore());

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({groupStore});
    });
    it('returns null if no group', () => {
        expect(shallow(<Group />)).toBeEmptyRender();
    });
    it('displays group name', () => {
        const id = 234;
        const name = 'the group';
        jest.spyOn(groupStore, 'getGroup').mockReturnValue({name} as GroupModel);

        const component = shallow(<Group id={id} />);

        expect(component.find('.group')).toHaveText(name);
        expect(groupStore.getGroup).toBeCalledWith(id);
    });
});
