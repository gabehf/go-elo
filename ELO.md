# How Elo is Calculated

## Calculating the expected chance to win

The probability a player has to win is calculated as follows:

$R_1=10^{Elo_1/V}$

$R_2=10^{Elo_2/V}$

$E_1=R_1/(R_1+R_2)$

$E_2=R_2/(R_1+R_2)$

Where:

- $Elo_1$ is the elo of player one,
- $Elo_2$ is the elo of player two,
- $V$ is the deviation parameter,
- $E_1$ is the probabilty of player one winning, and
- $E_2$ is the probability of player two winning.

There are two factors that influence the probability that a player has to win: the difference in elo, and the deviation parameter. The greater the difference in elo, the more likely it is for the higher-elo player to win, whereas the greater the deviation, the less likely the higher elo player to win.

## Unscored (win/loss) matches

Scored matches have two factors that determines how much the players' elo will change after the match. These variables are the chance the winner had to win, and the other is the $K$-Factor.

The higher the $K$-Factor, the more rapid changes in elo will occur. Basically, with a higher $K$-Factor, more elo will be gained per win, and more will be lost per loss.

## Scored matches

Scored matches have two variables that determine how much the players' elo will change after the match, other than the $K$-Factor and the probability the winner had to win.

- How dominant the scoreline was ($D$)
- The Score Weight parameter ($W$)

The greater the dominance factor $D$ is, the more elo will move around. A dominant victory will result in more elo being
given to the winner, and more being taken from the loser.

The greater the Score Weight parameter $W$ is, the more $D$ is taken into account when calculating elo, and the less elo
will move around in general. Increasing $W$ will result in only more dominant matches being given a large elo change, and closer
matches only moving a small amount of elo. This overall reduction in elo changes can be counteracted by increasing $K$, which will
result in closer matches being given a moderate elo change, and dominant matches moving a lot of elo.

The amount of elo to be lost or gained based on these factors can be written as the equation:

$S_1=((E_LD)e^{-WE_W}+E_W)$

$S_2=E_L-(E_LD)e^{-WE_W}$

$Elo_W=K(S_1-E_W)$

$Elo_L=K(S_2-E_L)$

Where:

- $E_W$ is the expected win chance of the winner,
- $E_L$ is the expected win chance of the loser,
- $D$ is calculated as $Score_W/{(Score_W+Score_L)}$
- $Elo_W$ is the amount of elo gained by the winner, and
- $Elo_L$ is the amount of elo gained by the loser (a negative value).

The $Elo$ values are then added to the players' current elo.
