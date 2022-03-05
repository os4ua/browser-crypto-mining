import Head from 'next/head'
import styles from '../styles/Home.module.scss'
import {ReactComponentElement} from "react";

interface Counter {
    count: number;
    name: string;
    icon: string;
}

export default function Home(): ReactComponentElement<any> {

  const counters: Counter[] = [
      {count: 1000000, name: "Users live", icon: "group"},
      {count: 50000, name: "Unique visitors", icon: "person"},
      {count: 10000000, name: "Hours spent on site", icon: "schedule"}
  ];

  return (
    <div className={styles.container}>
      <Head>
        <title>Browser Crypto Mining - OS4UA</title>
        <meta name="description" content="Browser Crypto Mining to help Ukraine" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <header className={styles.header}>
        <span>PukFutin</span>
        <nav className={styles.nav}>
          <a href="#">Why help Ukraine?</a>
          <a href="#">How it works</a>
          <a href="#">About donations</a>
          <a href="#">Some history</a>
        </nav>
        <div>
          <span>English</span>
        </div>
      </header>

      <main className={styles.main}>
    <h1>
      By being on this website for 00:00:00
      you have donated so far € 0,0.0
    </h1>
      <div className={styles.donations}>
          <span>Total donations so far</span>
          <p>€ 30,000</p>
          <div className={styles. bitcoinDonations}>
              <h4>Bitcoin wallet address</h4>
              <div><span>1Awyd1QWR5gcfrn1UmL8dUBj2H1eVKtQhg</span></div>
          </div>
      </div>
      <span>Don’t want to keep this page open?  <a href="#"> donate here</a></span>
      </main>
    <section className={styles.counters}>
        {counters.map((c,i) => (
          <div key={i} className={styles.counter}>
              <div>
                  <span className="material-icons">{c.icon}</span>
              </div>
              <p>{c.count}</p>
              <h5>{c.name}</h5>
          </div>
        ))}
    </section>
    </div>
  )
}
