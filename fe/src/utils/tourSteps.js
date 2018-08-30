const STRINGS = {
  welcome_title: {
    en: 'Welcome to profitcoins.io',
    de: 'Willkommen bei Profitcoins.io'
  },
  welcome_body: {
    en: 'In this quick tutorial you will see, how Profitcoins.io generates your first profits.',
    de: 'In diesem kurzen Tutorial lernen Sie, wie Sie mit Profitcoins.io konstante Gewinne erzielen können.'
  },
  balance_title: {
    en: 'Account Balance',
    de: 'Account Balance'
  },
  balance_body: {
    en: 'Your account balance shows the total amount of Bitcoin that you have on your Profitcoins.io account.',
    de: 'Die Account Balance zeigt Ihnen Ihren gesamten Profitcoins.io Kontostand an.'
  },
  profit_title: {
    en: 'Your Profit',
    de: 'Your Profit'
  },
  profit_body: {
    en: 'Your profit shows the amount of Bitcoin that you have earned with Profitcoins.io so far.',
    de: 'Ihr Profit zeigt Ihnen, wieviele Bitcoin Sie durch Profitcoins.io bereits verdient haben.'
  },
  margin_title: {
    en: 'Margin',
    de: 'Margin'
  },
  margin_body: {
    en: 'Here you can see your total margin in percent.',
    de: 'Hier sehen Sie Ihren gesamten Gewinn in Prozent.'
  },
  leverage_title: {
    en: 'Leverage',
    de: 'Leverage'
  },
  leverage_body: {
    en: 'Leverage refers to the leverage effect of our capital on your investment. Leverage 1:30 means that we increase your investment by the factor thirty. Thus, you can earn 30-times higher profits.',
    de: 'Die Leverage ist der Faktor, mit dem wir Ihr Kapital hebeln. Leverage 1:30 bedeutet, dass wir Ihr Investment um das 30-fache aufstocken. So können Sie 30 mal höhere Gewinne erzielen.'
  },
  demo_title: {
    en: 'Demo',
    de: 'Demo'
  },
  demo_body: {
    en: 'Press the Start Demo button now to activate your free demo.',
    de: 'Drücken Sie jetzt auf Start Demo um die kostenlose Demo-Version zu testen.'
  },
  dashboard_title: {
    en: 'Dashboard',
    de: 'Dashboard'
  },
  dashboard_body: {
    en: 'Great, you have just activated your demo! Here you can see, how Profitcoins.io buys and sells cryptocurrencies to generate your profits.',
    de: 'Sehr gut, Sie haben Ihre Demo aktiviert! Hier sehen Sie, wie Profitcoins.io für Sie Kryptowährungen günstig kauft und teuer verkauft.'
  },
  invest_title: {
    en: 'Invest',
    de: 'Invest'
  },
  invest_body: {
    en: 'As soon as you are ready for real profit, press the Invest button and credit your account.',
    de: 'Sobald Sie bereit sind echte Gewinne zu erzielen, klicken Sie auf Invest und laden Sie Ihr Konto auf.'
  }
};

export default (lang = 'en') => {
  lang = lang.toLowerCase();
  return [
    {
      title: STRINGS.welcome_title[lang],
      text: STRINGS.welcome_body[lang],
      selector: 'header',
      position: 'right',
      style: {
        backgroundColor: '#f6f7f8',
        arrow: {
          display: 'none'
        },
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.balance_title[lang],
      text: STRINGS.balance_body[lang],
      selector: '.tour-balance-highlighter',
      style: {
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.profit_title[lang],
      text: STRINGS.profit_body[lang],
      selector: '.tour-profit-highlighter',
      style: {
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.margin_title[lang],
      text: STRINGS.margin_body[lang],
      selector: '.profit__circle',
      style: {
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.leverage_title[lang],
      text: STRINGS.leverage_body[lang],
      selector: '.tour-leverage-highlighter',
      style: {
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.demo_title[lang],
      text: STRINGS.demo_body[lang],
      allowClicksThruHole: true,
      selector: '.demo-btn',
      style: {
        footer: { display: 'none' },
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.dashboard_title[lang],
      text: STRINGS.dashboard_body[lang],
      selector: '.dashboard-card',
      style: {
        close: { display: 'none' }
      }
    },
    {
      title: STRINGS.invest_title[lang],
      text: STRINGS.invest_body[lang],
      selector: '.invest-btn',
      style: {
        close: { display: 'none' }
      }
    }
  ];
};
