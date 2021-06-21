# vote

CLI app for voting information. Makes use of Google's [Civic Information API](https://developers.google.com/civic-information).

### Build

```shell
go build
```

### Usage
```
❯ ./vote
Vote is a CLI app for upcoming elections

Usage:
  vote [command]

Available Commands:
  elections   List upcoming elections
  help        Help about any command
  reps        List representatives

Flags:
  -h, --help   help for vote

Use "vote [command] --help" for more information about a command.
```

### Example

#### Reps command
```
❯ ./vote reps
What's your address?: 111 Eighth Avenue

New York Mayor
• Bill de Blasio
• Democratic Party

New York City Comptroller
• Scott M. Stringer
• Democratic Party
• action@comptroller.nyc.gov

New York Public Advocate
• Jumaane D. Williams
• Democratic Party
• reception@advocate.nyc.gov

New York City Council Member
• Corey Johnson
• Democratic Party
• speakerjohnson@council.nyc.gov
```

#### Elections command

```
❯ ./vote elections
? Pick an upcoming election:
  VIP Test Election |  June 6, 2025
  ❯️  New York Municipal Election |  June 22, 2021
  House District 78 Primary Runoff Election |  June 22, 2021
  Pawtucket City Council District 5 Special Election |  July 6, 2021
↓ Los Angeles County - San Marino Unified School District Special Parcel Tax Election |  June 29, 2021
```
```
❯ ./vote elections
❯️️  New York Municipal Election |  June 22, 2021
What's your address?: 111 Eigth Avenue

YOUR EARLY VOTE SITE:
The Mall at Hudson Yards
20 Hudson Yards
New York, NY, 10001

POLLING HOURS:
Thu, Jun 17: 10 am - 8 pm
Fri, Jun 18: 7 am - 4 pm
Sat, Jun 19: 8 am - 5 pm
Sun, Jun 20: 8 am - 4 pm

YOUR VOTE DAY POLLING SITE:
IS 70
333 West 17 Street
New York, NY, 10011

POLLING HOURS:
Tue, Jun 22: 6 am - 9 pm
```