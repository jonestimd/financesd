import {MenuList, MenuItem, Icon} from '@material-ui/core';
import React from 'react';
import {translate} from '../../i18n/localize';

interface IChildMenuProps {
    children: React.ReactNode;
    name?: string;
    onBack: () => void;
}

const ChildMenu: React.FC<IChildMenuProps> = ({children, name, onBack}) => {
    return (
        <MenuList>
            <MenuItem onClick={onBack} divider><Icon>chevron_left</Icon>{translate("menu.back")}</MenuItem>
            {name && <MenuItem className='menu-title' disabled divider>{name}</MenuItem>}
            {children}
        </MenuList>
    );
};

export default ChildMenu;
