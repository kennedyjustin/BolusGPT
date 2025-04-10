# BolusGPT

BolusGPT is an OpenAI Custom GPT that calculates [bolus](https://en.wikipedia.org/wiki/Bolus_(medicine)) insulin doses via natural language. Users can use text or voice for prompts like the following:

> Bolus dose for an 8 oz steak, a cup of cooked broccoli, and a cup of brown rice.

> Change my insulin-carb-ratio to 1:6.

The calculation is done on an HTTP server and the API is exposed to the Custom GPT via an "Action" (function calling). The server integrates with Dexcom CGMs (via a port of [pydexcom](https://github.com/gagebenne/pydexcom/)) to get the user's real-time blood glucose level and trend, and also stores static settings like the user's insulin-to-carb ratio.

**Important Note**: Insulin dosing is under the purview of the FDA, so users are required to self-host the server. All of the resources required to build the Custom GPT, along with a [setup guide](./SETUP.md) is included in the repository.

## Demo

TODO: Video

## What is included?

### Server
1. Stores static settings (in a simple JSON file) required to calculate bolus doses .For example:
   1. Target blood glucose level
   1. Insulin-to-carb ratio
   1. Insulin sensitivity factor
1. Retrieves the user's real-time blood glucose level and trend via the Dexcom API.
1. Calculates bolus insulin doses via well-tested algorithm.

### Custom GPT Resources
1. OpenAPI spec ([`openapi.yaml`](./openapi.yaml)) included so the CustomGPT knows how to interact with the APIs.
1. Spec for the CustomGPT itself ([`SPEC.md`](./SPEC.md))
1. Nutrition database file so the CustomGPT retrieves accurate nutrition information.

## How does it work?

TODO: Diagram

A brief outline of how this works end-to-end:

1. The user follows the [setup guide](./SETUP.md) and creates their own BolusGPT in their OpenAI account.
1. On their preferred OpenAI app (iOS App, website, etc.), and using their preferred medium of interaction (text / voice), the user interacts with BolusGPT.
1. First, they onboard.
   1. BolusGPT asks for the user's insulin-to-carb ratio, target blood glucose level, etc.
   1. BolusGPT calls the `PATCH /me` API with the information, where it is stored in a JSON file on the server.
1. Next, they ask BolusGPT to dose their meal.
   1. BolusGPT asks what the user will eat, whether the user will soon be exercising, etc.
   1. BolusGPT references the nutritional information database file to collect grams of carbs, fiber, protein, etc.
   1. BolusGPT calls the `POST /dose` API with the information.
      1. On the server, the current blood glucose level and trend is received from the Dexcom API.
      1. Using the nutrition info, exercise info, current blood glucose info, and stored user info, the insulin bolus is calculated and returned.
  1. BolusGPT presents the bolus dose to the user.
1. Optionally, the user can confirm they will use this dose (or tell BolusGPT they will opt for a different dose).
   1. If confirmed, BolusGPT can call `PATCH /me` with the dose used, and when. This can be an input to the next dose calculation (used to calculate insulin-on-board (IOB)).

### Documentation

#### `/me`

Any of the following fields can be retrieved (via `GET`), or updated (via `PATCH`):

- `fiber_multiplier` - Adjustment factor for dietary fiber's effect on insulin needs. A value of `1` counts all fiber. A value of `0` subtracts all fiber from total carbs
- `sugar_alcohol_multiplier` - Adjustment factor for sugar alcohols' impact on blood sugar. A value of `1` counts all sugar alcohol. A value of `0` subtracts all sugar alcohol from total carbs
- `protein_multiplier` - Factor representing how protein contributes to insulin demand. A value of `1` counts all protein. A value of `0` counts none of the protein
- `carb_threshold_to_count_protein_under` - Carb threshold under which protein is counted for dosing. For example, when the value is `20`, if the calculated carbs is under `20` protein is calculated according to the multiplier.
- `insulin_to_carb_ratio` - Grams of carbs covered by one unit of insulin. A value of `5` specifies 1 unit of insulin to 5 grams of carbs (1:5)
- `target_blood_glucose_level_in_mg_dl` - Target blood glucose level in mg/dL.
- `insulin_sensitivity_factor` - Blood glucose drop expected per unit of insulin. A value of `20` means a drop of 20 mg/dL is expected for 1 unit of insulin.
- `last_bolus_time` - Time of the last insulin bolus.
- `last_bolus_units_of_insulin` - Units of insulin used in the last bolus.

#### `/dose`

Any of the following fields can be provided when asking for a dose calculation via `POST`.

- `total_grams_of_carbs` - Total grams carbohydrates in the meal.
- `grams_of_fiber` - Grams of dietary fiber in the meal.
- `grams_of_sugar_alcohol` - Grams of sugar alcohols in the meal.
- `grams_of_protein` - Grams of protein in the meal.
- `minutes_of_exercise` - Duration of exercise in minutes that will occur after the bolus.
- `exercise_intensity` - Intensity of exercise that will occur after the bolus (`none`, `low`, `medium`, `high`).

## How can I use this?

Please follow the [setup guide](./SETUP.md).

## Acknowledgements

1. Gage Benne for reverse engineering the Dexcom Share API: https://github.com/gagebenne/pydexcom
1. Gary Scheiner's book, [Think Like a Pancreas](https://www.amazon.com/Think-Like-Pancreas-Practical-Insulin-Completely/dp/0738215147)
1. Dr. Richard Bernstein's book, [Diabetes Solution](https://www.amazon.com/Dr-Bernsteins-Diabetes-Solution-Achieving/dp/0316182699)

## Todo

- GPT Spec
- GPT needs nutrition database file
- Test
- Demo video
- Productionize
  - Put API on `api` subdomain
  - Direct `bolusgpt.com` to Github URL
