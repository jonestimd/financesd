import React from 'react';
import classNames from 'classnames';
import {CompanyModel} from '../../model/account/CompanyModel';
import MainMenu from './MainMenu';
import CompaniesMenu from './CompaniesMenu';
import AccountsMenu from './AccountsMenu';

const useSlideMenu = () => {
    const [depth, setDepth] = React.useState(0);
    const [slide, setSlide] = React.useState<'next' | 'previous' | null>(null);
    return {
        depth, slide,
        selectPane: (dir: -1 | 1) => {
            setSlide(dir > 0 ? 'next' : 'previous');
            setDepth(depth + dir);
        },
        onTransitionEnd: () => setSlide(null),
    };
};

interface IProps {
    currentPage: string;
}

const PageMenu: React.FC<IProps> = ({currentPage}) => {
    const {depth, slide, selectPane, onTransitionEnd} = useSlideMenu();
    const [company, setCompany] = React.useState<CompanyModel | undefined>(undefined);
    const selectCompany = (company?: CompanyModel) => {
        setCompany(company);
        selectPane(1);
    };
    const panels: (React.ReactElement | null)[] = [
        <MainMenu currentPage={currentPage} setDepth={() => selectPane(1)}/>,
        <CompaniesMenu onBack={() => selectPane(-1)} selectCompany={selectCompany} />,
        company ? <AccountsMenu onBack={() => selectPane(-1)} company={company} /> : null,
    ];
    return (
        <div className='sidebar-menu'>
            <div className={classNames('menu', slide)} onTransitionEnd={onTransitionEnd}>
                {slide === 'next' ? panels[depth-1] : panels[depth]}
                {slide === 'next' ? panels[depth] : slide && panels[depth + 1]}
            </div>
        </div>);
};

export default PageMenu;
